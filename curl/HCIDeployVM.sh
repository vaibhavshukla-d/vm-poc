curl -X POST \
  'https://us-west.api.greenlake.hpe.com/virtualization/v1beta1/virtual-machines' \
  --header 'Accept: application/json' \
  --header 'Authorization: Bearer YOUR_TOKEN' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "vmConfig": {
        "name": "string",
        "acceptEula": true
    },
    "storageConfig": {
        "defaultDatastoreId": "string"
    }
}'
