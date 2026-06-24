#include "../include/utils.h"
#include <stdio.h>
#include <stdlib.h>

extern char **environ;

void listEnvs() {
  int i = 0;
  while (environ[i] != NULL) {
    printf("%s\n", environ[i]);
    i++;
  }
}

int main(int argc, char *argv[]) {
  Args args = {.help = 0,
               .addEnv = 0,
               .envName = NULL,
               .envValue = NULL,
               .action = EnvList};

  if (parseArgs(&args, argc, argv)) {
    exit(1);
  }

  printf("%d\n", args.action);

  if (args.help) {
    char *message = "Manage environment variables\n"
                    "  -h --help              print this help message\n"
                    "  -a [Name] [Value]      add env variable\n";
    printf("%s\n", message);
    return 0;
  }

  switch (args.action) {
  case EnvAdd:
    printf("Adding env\n");
    break;
  case EnvUpdate:
    printf("Updating env\n");
    break;
  case EnvRemove:
    printf("Removing env\n");
    break;
  case EnvList:
    listEnvs();
    break;
  }
  return EXIT_SUCCESS;
}
