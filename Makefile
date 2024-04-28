oc_ingest: ingest_tools/owncloud/ingest.go
	CGO_ENABLED=0 go build -o ./oc_ingest ingest_tools/owncloud/ingest.go