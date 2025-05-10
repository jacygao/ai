docker build -t gcr.io/electric-facet-306612/hello-world:latest .

gcloud auth configure-docker

docker push gcr.io/electric-facet-306612/hello-world:latest