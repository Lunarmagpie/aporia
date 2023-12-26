#include <stdio.h>
#include <unistd.h>
#include <grp.h>
#include <security/pam_appl.h>
#include <security/pam_misc.h>

char* index_string_array(char** arr, int index) {
  printf("%s", arr[index]);
  return arr[index];
}
