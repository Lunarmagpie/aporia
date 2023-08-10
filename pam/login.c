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
  for (i = 0; i < num_msg; i++) {
      char *username, *password;
      switch (msg[i]->msg_style) {
        case PAM_PROMPT_ECHO_ON:
            username = ((char **) appdata_ptr)[0];
            (*resp)[i].resp = strdup(username);
            break;
        case PAM_PROMPT_ECHO_OFF:
            password = ((char **) appdata_ptr)[1];
            (*resp)[i].resp = strdup(password);
            break;
        case PAM_ERROR_MSG:
            fprintf(stderr, "%s\n", msg[i]->msg);
            result = PAM_CONV_ERR;
            break;
        case PAM_TEXT_INFO:
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

struct pam_conv new_conv(const char *username, const char *password) {
  const char *data[2] = {username, password};
  struct pam_conv conv = {
    conv_callback,
    data,
  };
  return conv;
} 
