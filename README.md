# ONEPASS

ONEPASS helps users establish [Cloudflare Tunnel](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/) with [`cloudflared`](https://developers.cloudflare.com/cloudflare-one/applications/non-http/cloudflared-authentication/arbitrary-tcp/).

## Build-time arguments

- `CLOUDFLARED_VERSION`: The version of `cloudflared` should be used.
- `REMOTE`: The remote address(es), separate with comma(`,`) of more than one.
- `LOCAL`: The local address(es), separate with comma(`,`) of more than one.
