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
    dns_zone = var.dns_zone
    alb_dns_name = var.alb_dns_name
    alb_https = true
    alb_http = false 
    ecs_app_env_map = var.ecs_app_env_map
    ecs_app_env_map_secret = var.ecs_app_env_map_secret
    ecs_app_env_map_secret_keys = var.ecs_app_env_map_secret_keys
}
