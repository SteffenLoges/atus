<template>
  <v-container class="d-flex mt-4" style="max-width: 1200px">
    <div class="flex-grow-1">
      <Suspense>
        <router-view />

        <template #fallback>
          <v-row align="center" justify="center" class="fill-height">
            <v-progress-circular indeterminate color="primary" size="64" />
          </v-row>
        </template>
      </Suspense>
    </div>
  </v-container>
</template>
