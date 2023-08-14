#include <stdio.h>
#include <unistd.h>
#include <grp.h>
#include <security/pam_appl.h>
#include <security/pam_misc.h>

char* index_string_array(char** arr, int index) {
  printf("%s", arr[index]);
  return arr[index];
}

// Set pam env for logind support
void set_pam_env(pam_handle_t *handle) {
	pam_misc_setenv(handle, "XDG_SESSION_CLASS", "greeter", 0);
	pam_misc_setenv(handle, "XDG_SEAT", "seat0", 0);
	pam_misc_setenv(handle, "XDG_VTNR", "1", 0);
}
