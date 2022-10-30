# Designing the service

Requirements
 - Generate reusable QRCodes 
 - Repurpose qrcodes


## How does it work?

An incoming request with data about a resource that needs to be encoded is sent to the server.
A response with a qrcode png code is returned.


## Research points
 - Can we create a qrcode with logo's embded?
  - In this case we get an base image and a resource we want to embed, the resulting generated
    QRCode is custom
 - Can we have some expiration on these codes, i.e codes that expire after a day. [WIP]