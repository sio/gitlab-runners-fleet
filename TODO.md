# TODO list

## build/

- Create CI pipeline with GitHub actions to push updated images to S3 (on schedule)
- Figure out Docker-in-Docker CI problems
- Switch reverse proxy from Caddy to Nginx: Caddy is not available in Debian repos


## deploy/

- Package into a Docker container
- Deploy fleet manager to home server
- Look into generating a pre-signed URL for VM image on fleet manager.
  That would allow to make the S3 bucket private.
  Be careful: changing URL (GET params) would trigger tf to rebuild the image,
  which in turn could(?) trigger VM rebuilds.


## scale/

- Rewrite in Go: calculate scaling actions and populate tfvars
- Optional: use external data source in terraform to call scaler automatically


## Global

- Write minimal documentation for each stage
