#include "../include/utils.h"
#include <stdio.h>
#include <string.h>

int parseArgs(Args *args, int argc, char **argv) {
  if (argc == 1)
    return 0;

  for (int i = 0; i < argc; i++) {
    if (strcmp(argv[i], "-h") == 0 || strcmp(argv[i], "--help") == 0) {
      args->help = 1;
    } else if (strcmp(argv[i], "-a") == 0) {
      args->action = EnvAdd;
      if (i + 2 < argc) {
        args->envName = argv[i + 1];
        args->envValue = argv[i + 2];
      } else {
        fprintf(stderr, "Error: -a requires envName and envValue.\n");
        return 1;
      }

    } else if (strcmp(argv[i], "-u") == 0) {
      args->action = EnvUpdate;
      if (i + 2 < argc) {
        args->envName = argv[i + 1];
        args->envValue = argv[i + 2];
      } else {
        fprintf(stderr, "Error: -u requires envName and envValue.\n");
        return 1;
      }

    } else if (strcmp(argv[i], "-r") == 0) {
      args->action = EnvRemove;
      if (i + 1 < argc) {
        args->envName = argv[i + 1];
      } else {
        fprintf(stderr, "Error: -r requires EnvName.\n");
        return 1;
      }
    }
  }
  return 0;
}
