## Deploy a marketplace licenses checker as pod in Managed Kubernetes   

#### Docker commands

1. Build docker image:
    ```
    docker build -t sample-check-license .
    ```

2. Deploy docker image to [**Yandex Container Registry**](https://cloud.yandex.ru/docs/container-registry/).
    ```
    REGISTRY_ID=crp72iu38e4jkx0fdp8a
    docker tag sample-check-license cr.yandex/$REGISTRY_ID/sample-check-license:1.0.0
    docker push cr.yandex/$REGISTRY_ID/sample-check-license:1.0.0
    ```

3. Correct helm-chart values for your docker image.

    Change `image.repository` values in *helm/values.yaml* to your docker name,\
    like `cr.yandex/$REGISTRY_ID/sample-check-license`.

4. Test package in your cluster *(optional)*.

  - create service-account's key:
    ```
    yc iam key create --service-account-id <svc-id> --output key.json
    ```
  - install:
    ```
    helm install --set-file saKeySecretKey=./key.json --set-string clusterId="<cluster-id>" --set-string licenseId="licenseId" sample-check-license ./helm/sample-check-license/
    ```
  - find your service:
    ```
    kubectl get services sample-check-licensesample-check-license-service
    NAME                           TYPE           CLUSTER-IP     EXTERNAL-IP     PORT(S)          AGE
    sample-check-license-service   LoadBalancer   ...            ...             8080:32468/TCP   10
    ```
  - ping service:
    ```
    curl <your external ip>:8080
    License is ERROR: %!(EXTRA string=rpc error: code = Unavailable desc = Unavailable)
    ```
  - uninstall:
    ```
    helm uninstall sample-check-license
    ```

5. Push helm chart to your container registry:
    ```
    helm package ./helm/sample-check-license/
    helm push sample-check-license-1.0.tgz oci://cr.yandex/$REGISTRY_ID/sample-check-license
    ```

6. Change helm/manifest.yaml value `helm_chart.name` to helm chart name from your registry, like
    ```
    cr.yandex/$REGISTRY_ID/sample-check-license/sample-check-license:1.0
    ```
   
7. Now you can publish your product by [**publisher instruction**](https://cloud.yandex.ru/docs/marketplace/operations/create-container)
