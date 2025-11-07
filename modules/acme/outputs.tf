output "certificate" {
   value =  acme_certificate.cert.certificate_pem
}

output "private_key" {
   value = tls_private_key.cert_private_key.private_key_pem
}