FROM iconbridge/build as build
COPY . bmr
RUN cd bmr/cmd/iconbridge && go build
RUN mv /bmr/cmd/iconbridge/iconbridge /bin/iconbridge
