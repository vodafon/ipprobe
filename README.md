# IpProbe

Lookup IP tool

### Installation and Upgrading

```bash
go get -u github.com/vodafon/ipprobe
```

### Arguments

```
Usage of ipprobe:
  -format6
        URL friendly IPv6
  -host
        print host too
  -ipv4
        only IPv4
  -ipv6
        only IPv6
  -procs int
        concurrency (default 20)
```

## Usage

```bash
echo "google.com"|ipprobe

172.217.16.14
2a00:1450:401b:804::200e
```

### IPv4 only

```bash
echo "google.com"|ipprobe -ipv4

172.217.16.14
```

### IPv6 only

```bash
echo "google.com"|ipprobe -ipv6

2a00:1450:401b:805::200e
```

### Format IPv6 for URLs

```bash
echo "google.com"|ipprobe -format6

172.217.20.206
[2a00:1450:401b:804::200e]
```

### With host

```bash
cat hosts.txt

google.com
facebook.com
twitter.com
```

```bash
cat hosts.txt|ipprobe -format6 -host

104.244.42.193 twitter.com
104.244.42.65 twitter.com
172.217.20.206 google.com
[2a00:1450:401b:804::200e] google.com
31.13.81.36 facebook.com
[2a03:2880:f116:83:face:b00c:0:25de] facebook.com
```
