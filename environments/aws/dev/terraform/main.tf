module infra {
    source = "../../../../modules/aws/infra"
    env = var.env
    docker_image = var.docker_image
    sshpubkeypath = var.sshpubkeypath
    sshprivkeypath = var.sshprivkeypath
    region = var.region
    ssm_path = var.ssm_path
    bastion = false
    private_host = false
}
