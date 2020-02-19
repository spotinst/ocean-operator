FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

# set the environment variables
ENV OPERATOR_NAME=ocean-operator \
    USER_NAME=ocean-operator \
    USER_UID=1001

# set the metadata labels
LABEL vendor="Spotinst" \
      name="ocean-operator" \
      description="Ocean Operator for Kubernetes" \
      summary="Ocean Operator for Kubernetes is an Operator that makes use of custom resources and can be used to create and manage Ocean resources in Spot"

# copy licenses
COPY LICENSE /licenses/spotinst/LICENSE

# install operator binary
COPY build/_output/${OPERATOR_NAME}_linux_amd64/${OPERATOR_NAME} /usr/local/bin/${OPERATOR_NAME}

# install helpers
COPY build/bin /usr/local/bin

# set up the user
RUN /usr/local/bin/user_setup

# set the current user
USER ${USER_UID}

# set the entrypoint
ENTRYPOINT ["/usr/local/bin/entrypoint"]
