#include "../include/utils.h"
#include <stdio.h>
#include <stdlib.h>

extern char **environ;

int main(int argc, char *argv[]) {
  Args args = {.help = 0,
               .addEnv = 0,
               .envName = NULL,
               .envValue = NULL,};

  if (parseArgs(&args, argc, argv)) {
    exit(1);
  }

  if (args.help) {
    char *message = "Manage environment variables\n"
                    "  -h --help              print this help message\n"
                    "  -a [Name] [Value]      add env variable\n";
    printf("%s\n", message);
    return 0;
  } else if (args.addEnv) {

  } else {
    int i = 0;
    while (environ[i] != NULL) {
      printf("%s\n", environ[i]);
      i++;
    }
  }
  return EXIT_SUCCESS;
}
