docker build --no-cache -t gritter_backend_prod_linux --platform linux/amd64 .
docker tag gritter_backend_prod_linux asia-east1-docker.pkg.dev/gritter/backend/prod
docker push asia-east1-docker.pkg.dev/gritter/backend/prod
