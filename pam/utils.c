#include <stdio.h>

char* index_string_array(char** arr, int index) {
  printf("%s", arr[index]);
  return arr[index];
}
