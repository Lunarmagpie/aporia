#include <security/pam_appl.h>
#include <security/pam_misc.h>

static int conv_callback(
  int num_msg,
  const struct pam_message **msg,
  struct pam_response **resp,
  void *appdata_ptr
) {
    int i;

    *resp = calloc(num_msg, sizeof(struct pam_response));
    if (*resp == NULL) {
        return PAM_BUF_ERR;
    }

    int result = PAM_SUCCESS;

    const char *password = appdata_ptr;

    for (i = 0; i < num_msg; i++) {
        switch (msg[i]->msg_style) {
            case PAM_PROMPT_ECHO_ON:
                // This is called when asking for the username.
                // That will never happen.
                break;
            case PAM_PROMPT_ECHO_OFF:
                resp[i]->resp = strdup(password);
                break;
            case PAM_ERROR_MSG:
                // TODO
                fprintf(stderr, "%s\n", msg[i]->msg);
                result = PAM_CONV_ERR;
                break;
            case PAM_TEXT_INFO:
                // TODO
                printf("%s\n", msg[i]->msg);
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

struct pam_response* allocate_conv(int num_msgs) {
    return calloc(num_msgs, sizeof(struct pam_response));
}

struct pam_conv new_conv(char *password) {
    struct pam_conv conv = {
        conv_callback,
        password,
    };
    return conv;
} 
