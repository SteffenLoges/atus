<template>
  <Card title="Authorisation" variant="text" class="card-accent mb-4">
    <v-card-text>
      All requests to the API must be authenticated with the API Authentification token above.
      <br>
      <br>
      Make sure to include the token either as a header with the name <code>Authorization</code> or as a
      query parameter with the name <code>token</code> in every request.
    </v-card-text>
  </Card>
</template>