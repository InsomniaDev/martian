FROM nginx:1.13.3

LABEL maintainer="InsomniaDev <insomniadev@gmail.com>"

ARG REACT_APP_KEYCLOAK_REALM_NAME

COPY web-server-config/nginx-default.conf /etc/nginx/conf.d/default.conf
COPY build /usr/share/nginx/html

EXPOSE 80
