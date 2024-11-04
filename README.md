# United AI Prototypes

> this repo is a showcase of various prototypes for United AI

## License

i have released this source texts under [Apache 2.0 License](LICENSE.md) for
easy collaboration and let United AI easily move forward if source authors
are no longer a member.

## GCP Prototype

> hosted in GCP

**check out the demo at [uai-gcp.methebe.com](https://uai-gcp.methebe.com/)**

the `gcp/` folder contains a prototype of golang app which can be used to
elect project members. the aim is to experiment with a server-side binary that
is portable to any cloud provider or on-prem services.

this have no vendor lock-in whatsoever and doesn't need a dedicated DB server.
since usecases are pretty basic, this app uses _sqlite_. the frontend is
a golang templates inside `gcp/views/`.

### Build

to build this project, a script is provided (`gcp/build.sh`). the build server
should have go runtime and developer environment, node.js (with pnpm) and upx
for binary reduction. along with these, the build server also needs an active
internet connection for the first-time.

> build script doesn't manage build endvironment. it is your responsibility to
> ensure that the build script have access to these. **building user should be
> able to wrtie to filesystem in-place**

after the script is done, following files should be placed as it is in the
tree to deployment envrionment:

- uai (binary)
- style.css
- views/\*
- uai.db (optional)

### Features

- portable
- datastore is in-built
- easy to maintain
- can be entended to host main landing page as well
- begineer-friendly, so easy to find new maintainers
- can be containerized for easy deployments

### Drawbacks

- requires VPS with a presistent disk

## Vercel Prototype

> hosted in vercel

**check out the demo at [uai-vercel.methebe.com](https://uai-vercel.methebe.com/)**

this is just an astro project in SSR mode. supabase is
utilized as third-party datastore since vercel doesn't
provide persistent disk.

this prototype uses pqsql for obvious reasons.

### Features

- not portable as astro ssr adapter for vercel is unique
- relatively beginner friendly than gcp

### Drawbacks

- might end up in vendor-lock-in
- expensive option if free tier is exceeded

## Other Providers

gcp prototype can be deployed to any VPS provider such as render,
digital ocean etc.

vercel prototype can be deployed to netlify,
cloudflare edge etc. where astro supports SSR mode.

---

Servus!

Birnadin Erick.

Made is Deggendorf.
