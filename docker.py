import yaml
import os
import argparse
import itertools

class str2list():
    def __init__(self, out_type, delimiter=","):
        self.type = out_type
        self.delimiter = delimiter

    def __call__(self, value):
        values = []

        if value is not None:
            values = list(map(lambda str: self.type(str), value.split(self.delimiter)))

        return values

def setup_arguments():
    parser = argparse.ArgumentParser(description="Generating peers docker file for Promark.")
    parser.add_argument("-p", "--peers", type=int, required=True)
    parser.add_argument("-o", "--orgs", type=int, required=True)
    parser.add_argument("-t", "--types", type=str2list(str), default="adv,pub")

    args = parser.parse_args()

    params = {}

    for arg in vars(args):
        params[arg] = getattr(args, arg)

    return params

def get_peers_template_path():
    return os.path.join("docker", "template", "peer.yml")

def get_peers_dockerfile_path(num_orgs_per_type, num_peers_per_org):
    return os.path.join("docker", "peers-{}-{}.yml".format(num_orgs_per_type, num_peers_per_org))

def generate_peers_docker_file(org_types, num_orgs_per_type, num_peers_per_org):
    data = {
        "version": "2",
        "networks": {
            "promarknet": {
                "external": True,
                "name": "promarknet"
            }
        }
    }

    services_data = {}

    with open(get_peers_template_path(), "r") as f:
        peer_template = f.read()

    # print(peer_template)

    base_port = {
        "adv": 5000,
        "pub": 6000,
    }
    peer_port_step=10
    org_port_step=100

    for org_type, org_idx, peer_idx in itertools.product(org_types, range(num_orgs_per_type), range(num_peers_per_org)):
        peer_name = "peer{}.{}{}.promark.com".format(peer_idx, org_type, org_idx)
        # peerPort=$((baseAdvPort + orgId * orgPortStep + peerPortStep * peerId))
        # dbPort=$((peerPort + 2))
        # apiPort=$((peerPort + 1))
        # peerName="peer${peerId}.${orgName}"
        peer_port = base_port[org_type] + org_idx * org_port_step + peer_port_step * peer_idx
        api_port = peer_port + 2
        db_port = peer_port + 1

        org_name = "{}{}".format(org_type, org_idx)
        current_peer_data = peer_template.replace("${PEER_ID}", str(peer_idx))
        current_peer_data = current_peer_data.replace("${ORG_NAME}", org_name)
        current_peer_data = current_peer_data.replace("${PROJECT_NAME}", "promark")
        current_peer_data = current_peer_data.replace("${PEER_PORT}", str(peer_port))
        current_peer_data = current_peer_data.replace("${API_PORT}", str(api_port))
        current_peer_data = current_peer_data.replace("${DB_PORT}", str(db_port))

        # print(current_peer_data)
        yaml_peer_data = yaml.safe_load(current_peer_data)
        # print(yaml_peer_data.keys())

        services_data.update(yaml_peer_data)
    # print(services_data.keys())
    # raise Exception()

    #     print(org_type)
    # for org_type in org_types:
    #     for org_idx in range(num_orgs_per_type):
    #         for peer_idx in range(num_peers_per_org):
    #             pass

    data["services"] = services_data

    docker_file_path = get_peers_dockerfile_path(num_orgs_per_type, num_peers_per_org)
    with open(docker_file_path, "w") as f:
        yaml.safe_dump(data, f)
        print("saved dockerfile at:{}".format(docker_file_path))


def main(args):
    num_peers_per_org = args["peers"]
    num_orgs_per_type = args["orgs"]
    org_types = args["types"]
    print("Generating dockerfile with args: {}".format(args))
    # read peers.yml

    generate_peers_docker_file(org_types, num_orgs_per_type, num_peers_per_org)
    # path = get_peers_template()
    # with open(path, "r") as f:
    #     data = yaml.safe_load(f)
    #     print(data)

if __name__ == "__main__":
    args = setup_arguments()
    main(args)