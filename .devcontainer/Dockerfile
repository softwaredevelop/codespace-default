ARG VARIANT=latest
FROM docker.io/library/ubuntu:$VARIANT

ARG USERNAME=codespace
ARG USER_UID=1000
ARG USER_GID=$USER_UID

ENV DEBIAN_FRONTEND=noninteractive

COPY library-scripts/* /tmp/scripts/

RUN \
    bash /tmp/scripts/common.sh "true" "debian"; \
    bash /tmp/scripts/setup-user.sh "${USERNAME}" "${USER_UID}" "${USER_GID}" "debian"; \
    bash /tmp/scripts/hadolint.sh "2.12.0" "debian"; \
    bash /tmp/scripts/shfmt.sh "3.6.0" "debian"; \
    bash /tmp/scripts/shellcheck.sh "0.9.0" "debian"; \
    bash /tmp/scripts/editorconfig.sh "2.6.0" "debian"; \
    bash /tmp/scripts/trivy.sh "0.36.0" "debian"; \
    apt-get autoremove; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/* /tmp/scripts
