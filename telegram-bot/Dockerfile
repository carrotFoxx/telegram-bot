FROM python:3.8

ENV PYTHONPATH "/opt/app:${PYTHONPATH}"
EXPOSE 8080

ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && dpkg-reconfigure -f noninteractive tzdata

CMD ["python3", "/opt/app/main.py"]
# dependencies
COPY requirements /opt/app/
RUN pip install --no-cache-dir -r /opt/app/requirements
# version info
ARG PLATFORM_VERSION=unknown
ENV PLATFORM_VERSION $PLATFORM_VERSION
# sources
COPY src /opt/app