struct pam_conv_appdata { char *password; };

struct pam_conv new_conv(
    char *password
);
