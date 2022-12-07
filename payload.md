### Create QrCode 
	method: POST
	payload_body: {
		resource: "string"
		site_id: "string"
	}

### Get QrCode
	method: GET
	payload: {
		id: "string"
		site_id: "string"
	}

### Update QrCode
	method: POST
	payload_body: {
		id: "string"
		resource "string"
		site_id: "string"
	}

### Delete QrCode
	method: POST
	payload_body: {
		id: "string"
		site_id: "string"
	}
