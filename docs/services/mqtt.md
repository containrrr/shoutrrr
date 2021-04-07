# MQTT

## URL Format

_mqtt://**`host`**:**`port`**?topic=**`topic`**_

## Optional parameters

You can optionally specify the **`disableTLS`**, **`clientId`**, **`username`** and **`password`** parameters in the URL:  
_mqtt://**`host`**:**`port`**?topic=**`topic`**&disableTLS=true&clientId=**`clientId`**&username=**`username`**&password:**`password`**_

## Parameters Description

- **Host** - MQTT broker server hostname or IP address (**Required**)  
  Default: _empty_  
  Aliases: `host`

- **Port** - MQTT server port, common ones are 8883 for TCP/TLS and 1883 for TCP (**Required**)
  Default: `8883`

- **Topic** - Topic where the message is sent (**Required**)
  Default: _empty_  
  Aliases: `Topic`

- **DisableTLS** - disable TLS/SSL Configurations  
  Default: `false`

- **ClientID** - The client identifier (ClientId) identifies each MQTT client that connects to an MQTT  
  Default: _empty_
  Aliases: `clientId`

- **Username** - name of the sender to auth  
  Default: _empty_
  Aliases: `clientId`

- **Password** - authentication password or hash  
  Default: _empty_  
  Aliases: `password`

## Certificates to use TCP/TLS

To use TCP/TLS connection, it is necessary the files:

- Cerficate Authority: ca.crt
- Client Certificate: client.crt
- Client Key: client.key

## Configure TLS in mosquitto

Generate the certificates [mosquitto-tls](https://mosquitto.org/man/mosquitto-tls-7.html).
