FROM envoyproxy/envoy:v1.17.0

COPY ./.docker/envoy/envoy.yaml /tmpl/envoy.yaml.tmpl
COPY ./.docker/envoy/docker-entrypoint.sh .

RUN chmod 500 /docker-entrypoint.sh

RUN apt-get update && apt-get install gettext -y

EXPOSE 8090
EXPOSE 9902

ENTRYPOINT ["./docker-entrypoint.sh"]
