FROM golang:1.12

RUN mkdir -p /workspace
WORKDIR /workspace

ENV KUSTOMIZE_TAG v3.1.0

RUN mkdir -p sigs.k8s.io && \
  git clone https://github.com/kubernetes-sigs/kustomize.git sigs.k8s.io/kustomize && \
  cd sigs.k8s.io/kustomize && \
  git checkout tags/$KUSTOMIZE_TAG -b $KUSTOMIZE_TAG && \
  go build -o /bin/kustomize cmd/kustomize/main.go && \
  mkdir -p ~/.config/kustomize/plugin/ryane/v1beta1/kfilter/

COPY kfilt.go sigs.k8s.io/kustomize/plugin/

RUN cd sigs.k8s.io/kustomize && \
  go get github.com/ryane/kfilt && \
  go build -buildmode plugin -o ~/.config/kustomize/plugin/ryane/v1beta1/kfilter/Kfilter.so plugin/kfilt.go

ENTRYPOINT ["/bin/kustomize"]
