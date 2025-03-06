FROM gcr.io/distroless/base-debian12
ARG TARGETARCH

WORKDIR /app

COPY dist/.env ./
COPY dist/${TARGETARCH}/satellite ./

EXPOSE 4000

ENTRYPOINT ["./satellite"]
CMD ["web", "serve"]
