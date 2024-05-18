#include <security/pam_appl.h>
#include <security/pam_misc.h>

// imports <login.h>
#include <_cgo_export.h>

static int conv_callback(
  int num_msg,
  const struct pam_message **msg,
  struct pam_response **resp,
  void* appdata_ptr
) {
    int i;

    *resp = calloc(num_msg, sizeof(struct pam_response));
    if (*resp == NULL) {
        return PAM_BUF_ERR;
    }

    int result = PAM_SUCCESS;

    struct pam_conv_appdata *appdata = appdata_ptr;
    const char *password = appdata->password;

    for (i = 0; i < num_msg; i++) {
        switch (msg[i]->msg_style) {
            case PAM_PROMPT_ECHO_ON:
                // This is called when asking for the username.
                // That will never happen.
                break;
            case PAM_PROMPT_ECHO_OFF:
                (*resp)[i].resp = strdup(password);
                break;
            case PAM_ERROR_MSG:
                handlePamErrorMessage((char*) msg[i]->msg);
                result = PAM_CONV_ERR;
                break;
            case PAM_TEXT_INFO:
                handlePamTextInfo((char*) msg[i]->msg);
                break;
        }
        if (result != PAM_SUCCESS) {
            break;
        }
    }

    if (result != PAM_SUCCESS) {
        free(*resp);
        *resp = 0;
    }

    return result;
}

struct pam_conv new_conv(char *password) {
    struct pam_conv_appdata *appdata = malloc(sizeof(struct pam_conv_appdata));

    appdata->password = password;
    
    struct pam_conv conv = {
        conv_callback,
        appdata,
    };

    return conv;
} 
