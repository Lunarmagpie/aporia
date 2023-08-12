#include <stdio.h>
#include <unistd.h>
#include <grp.h>

char* index_string_array(char** arr, int index) {
  printf("%s", arr[index]);
  return arr[index];
}
