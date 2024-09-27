# vcluster operator

vcluster operator is a Kubernetes operator that allows you to run multiple Kubernetes clusters on a single Kubernetes cluster.

## Usage
```bash
make docker # build the operator docker image
make deploy # deploy the operator to the current cluster
make forward # forward the operator to localhost:8080
make build # build the operator cli
export VCLUSTER_OPERATOR_URL=http://localhost:8080 # set the operator url
./bin/vcluster-operator --help # show the operator cli help
./bin/vcluster-operator login -U root -P admin # login to the operator
./bin/vcluster-operator create test # create a new vcluster
./bin/vcluster-operator list # list all current vclusters
./bin/vcluster-operator kubeconfig test # get the kubeconfig for a vcluster
./bin/vcluster-operator delete test # delete a vcluster
./bin/vcluster-operator logout # logout from the operator
```

## TODO:
- [ ] Add more documentation
- [x] Implement list
- [x] Implement create
- [x] Implement delete
- [ ] Integrate vault to create ephemeral credentials to be used in the pipeline
- [ ] make it curl friendly (remove dependency from httpie) (httpie is only needed for /login)
- [x] write a small client
- [ ] encode allowed cluster in the jwt
- [ ] implement auto delete (add labels to cluster)
- [x] fix /vcluster/[name]/kubeconfig
- [ ] print errors in client
- [ ] dump deployment to yaml
- [ ] add possibility to add certificates etc. to deployment

# Vault dynamic secrets provider

## Usage

All commands can be run using the provided [Makefile](./Makefile). However, it may be instructive to look at the commands to gain a greater understanding of how Vault registers plugins. Using the Makefile will result in running the Vault server in `dev` mode. Do not run Vault in `dev` mode in production. The `dev` server allows you to configure the plugin directory as a flag, and automatically registers plugin binaries in that directory. In production, plugin binaries must be manually registered.

This will build the plugin binary and start the Vault dev server:

```bash
# Build Vcluster plugin and start Vault dev server with plugin automatically registered
$ make vault
```

Now open a new terminal window and run the following commands:

```bash
# Open a new terminal window and export Vault dev server http address
$ export VAULT_ADDR='http://127.0.0.1:8200'

# Enable the Vcluster plugin
$ make vault-enable

# Write a secret to the Vcluster secrets engine
export VAULT_CLIENT_TIMEOUT=12000s
$ vault write vcluster-operator/config url="http://localhost:8080" username="admin" password="admin"
Success! Data written to: vcluster-operator/config

# Retrieve secret from Vcluster secrets engine
$ vault read vcluster-operator/test
Key      Value
---      -----
jwt      eJ...
```

## TODO
- [ ] implement config writer to put vcluster url and vcluster admin credentials
