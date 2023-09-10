# DHT Infrastructure layer

This is the foundational layer of our DHT. This includes the underlying
networking infrastructure that implements the protocols to `STORE` and `QUERY`
items into our network. This layer is very application agnostic. By design, we
support querying instead of lookup since our primary focus was to build a DHT
that can return not just exact matches but a ranges response including multiple
results that are semantically close to our search query. We can further tune the
search parameters to get exact matches, if that suits the application.

## Distributed Hash Table

- what kind DHT is it that we want to use
- we parametrize the hash function needed to compute the node ID
- In the default implementation, what will be our node ID

### Prefix Hash tree

## Content

hash of the content, content itself (for consensus only), url or node where it
is stored(to be handled by the application-interface layer)

Finger table of each node will have:

- <nodeID, IP, prefix>
