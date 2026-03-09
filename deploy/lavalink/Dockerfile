FROM ghcr.io/lavalink-devs/lavalink:4

COPY deploy/lavalink/application.yml /opt/Lavalink/application.yml

ENTRYPOINT ["java", "-Djdk.tls.client.protocols=TLSv1.2", "-jar", "/opt/Lavalink/Lavalink.jar"]
