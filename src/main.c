#include <stdio.h>
#include <stdlib.h>

extern char **environ;

int main(int argc, char *argv[]) {
  printf("Hello WayENV\n");
  int i = 0;

  while (environ[i] != NULL) {
    fprintf(stdout, "%s\n", environ[i]);
    i++;
  }
  return EXIT_SUCCESS;
}
