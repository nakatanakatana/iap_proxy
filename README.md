# iap_proxy

## How to use

1. Add `Authorized redirect URIs` to `http://localhost:18000/__/redirect` in OAuth Client for Cloud IAP
   * https://console.cloud.google.com/apis/credentials/
2. run iap_proxy
   * `docker run -e "IAP_PROXY_CLIENT_ID=${CLIENT_ID}" -e "IAP_PROXY_CLIENT_SECRET=${CLIENT_SECRET}" -e "IAP_PROXY_BASE_URL=${BACKEND_URL}" -p 18000:18000 ghrc.io/nakatanakatana/iap_proxy:latest`
3. access `http://localhost:18000/__/login` and login
