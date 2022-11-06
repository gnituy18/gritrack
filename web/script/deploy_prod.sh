docker build --no-cache -t gritrack_web_prod_linux --platform linux/amd64 .
docker tag gritrack_web_prod_linux asia-east1-docker.pkg.dev/gritrack/web/prod
docker push asia-east1-docker.pkg.dev/gritrack/web/prod
