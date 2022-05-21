from tikv_client import RawClient

client = RawClient.connect("127.0.0.1:2379")

# scan
print(client.scan(b":/", end=b":/", limit=100, include_start=True, include_end=True))

# client.delete_range(b"", b"", )