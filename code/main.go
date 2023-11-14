package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/marketplace/licensemanager/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/iamkey"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	clusterID  string
	licenseID  string
	secretFile string
	uuid       string
	port       int
	endpoint   string

	flags = registerFlags()
	sdk   *ycsdk.SDK
)

func main() {
	err := flags.Parse(os.Args[1:])
	if err != nil || len(clusterID) == 0 || len(licenseID) == 0 || len(secretFile) == 0 || len(uuid) == 0 {
		flags.PrintDefaults()
		os.Exit(-1)
	}

	sdk, err = buildSDK()
	if err != nil || sdk == nil {
		os.Exit(-2)
	}

	http.HandleFunc("/", handler)
	err = http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		os.Exit(-3)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	err := checkLicense(context.Background(), sdk, licenseID, clusterID+"_"+uuid)
	if err == nil {
		fmt.Fprintf(w, "License is OK")
		return
	}
	fmt.Fprintf(w, "License is ERROR: ", err.Error())
}

func checkLicense(ctx context.Context, sdk *ycsdk.SDK, licenseInstanceID string, resourceID string) error {
	ensureOp, err := sdk.Marketplace().LicenseManager().Lock().Ensure(ctx, &licensemanager.EnsureLockRequest{
		InstanceId: licenseInstanceID,
		ResourceId: resourceID, // Use compute instance id as resource ID
	})
	op, err := sdk.WrapOperation(ensureOp, err)
	if err != nil {
		return err
	}
	err = op.Wait(ctx)
	if err != nil {
		return err
	}
	if opErr := op.Error(); opErr != nil {
		return opErr
	}
	resp, err := op.Response()
	if err != nil {
		return err
	}
	lock := resp.(*licensemanager.Lock)
	fmt.Println("LockID:", lock.GetId())
	fmt.Println("Start:", lock.GetStartTime().AsTime())
	fmt.Println("End:", lock.GetEndTime().AsTime())
	return err
}

func registerFlags() *flag.FlagSet {
	f := &flag.FlagSet{}
	f.StringVar(&clusterID, "cluster-id", "", "cluster id")
	f.StringVar(&licenseID, "license-id", "", "license id")
	f.StringVar(&secretFile, "secret-file", "", "secret file")
	f.StringVar(&uuid, "uuid", "", "uuid")
	f.IntVar(&port, "port", 8080, "port")
	flag.StringVar(&endpoint, "endpoint", "", "cloud environment endpoint (defaults to prod endpoint)")
	return f
}

func buildSDK() (*ycsdk.SDK, error) {
	key, err := getCredsFromFile(secretFile)
	if err != nil {
		return nil, err
	}
	var creds ycsdk.Credentials
	creds, err = ycsdk.ServiceAccountKey(key)
	if err != nil {
		return nil, err
	}
	return ycsdk.Build(context.Background(), ycsdk.Config{
		Credentials: creds,
		Endpoint:    endpoint,
	})
}

type Key struct {
	Id               string `json:"id"`
	PrivateKey       string `json:"private_key"`
	ServiceAccountId string `json:"service_account_id"`
}

func getCredsFromFile(keyFile string) (*iamkey.Key, error) {
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	key := &Key{}
	err = json.Unmarshal(data, key)
	if err != nil {
		return nil, err
	}
	return &iamkey.Key{
		Id:         key.Id,
		Subject:    &iamkey.Key_ServiceAccountId{ServiceAccountId: key.ServiceAccountId},
		PrivateKey: key.PrivateKey,
	}, nil
}