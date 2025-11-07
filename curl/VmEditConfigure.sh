curl -i -X POST http://ind-south.api.qa-greenlake.hpe.com/virtualization/v1beta1/virtual-machines/101/update-hardware \
  --header 'Authorization: Bearer TEST' \
  --header 'Content-Type: application/json' \
  --data '{
    "cpuMemConfig": {
      "numOfCoresPerSocket": 4,
      "numOfCpus": 8,
      "memory": {
        "memoryInMb": 8192
      }
    },
    "networkAdapters": [
      {
        "connectAtPowerOn": true,
        "name": "eth0",
        "networkDetails": {
          "name": "VM Network",
          "type": "STANDARD_PORT_GROUP"
        },
        "operation": "ADD",
        "type": "VMXNET3"
      }
    ],
    "virtualDisks": [
      {
        "diskConfig": {
          "capacityInMb": 20480,
          "id": "123e4567-e89b-12d3-a456-426614174000",
          "retainFiles": false,
          "type": "SATA"
        },
        "operation": "ADD"
      }
    ]
  }'