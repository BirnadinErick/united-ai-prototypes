# United AI Prototypes

> this repo is a showcase of various prototypes for United AI

## License

i have released this source texts under [Apache 2.0 License](LICENSE.md) for
easy collaboration and let United AI easily move forward if source authors
are no longer a member.

## GCP Prototype

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

---

Servus!

Birnadin Erick.

Made is Deggendorf.
