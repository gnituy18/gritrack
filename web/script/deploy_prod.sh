docker build --no-cache -t gritter_web_prod_linux --platform linux/amd64 .
docker tag gritter_web_prod_linux asia-east1-docker.pkg.dev/gritter/web/prod
docker push asia-east1-docker.pkg.dev/gritter/web/prod
