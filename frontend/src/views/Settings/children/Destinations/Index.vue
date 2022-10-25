<template>
  <ConfirmDialog :modelValue="uploadSuccessMessage !== null" @always="uploadSuccessMessage = null"
    @confirm="($refs.form as HTMLFormElement).submit()" title="Upload successful">
    <v-alert type="success" variant="text" density="compact">The file was uploaded successfully.</v-alert>

    <v-alert type="warning" variant="tonal" density="compact" class="mt-2 mb-4">
      Please check the uploaded file on your tracker.<br />
      Try to download the torrent file and check if it is valid by adding it to a torrent client.
    </v-alert>

    <!-- <small
      >The Server responded with the following
      message:</small
    >
    <v-alert
      type="info"
      variant="tonal"
      border="start"
      density="compact"
      :icon="false"
    >
      {{ uploadSuccessMessage }}
    </v-alert> -->

    <p>Do you want to save these settings?</p>
  </ConfirmDialog>

  <ErrorDialog :modelValue="uploadErrorMessage !== null" @dismiss="uploadErrorMessage = null" title="Upload failed">
    <v-alert type="error" variant="text" density="compact" class="mb-3">
      The file could not be uploaded.<br />
      Check your settings and try again.
    </v-alert>

    <small>The Server responded with the following message:</small>
    <v-alert type="info" variant="tonal" border="start" density="compact" :icon="false">
      <span class="text-grey" style="white-space: pre-wrap">{{ uploadErrorMessage }}</span>
    </v-alert>
  </ErrorDialog>

  <FormCard ref="form" :loading="isLoading" title="Upload Settings" @submit="onSubmit">
    <v-card-text>
      <v-alert type="info" class="mb-4">
        Make sure you have installed and configured the
        <a href="https://github.com/SteffenLoges/atus-tracker-api" target="_blank" v-text="'ATUS Tracker API'"></a>
        on your tracker.
      </v-alert>

      <v-card variant="text" title="Main Settings" class="card-accent mb-4">
        <v-card-text>
          <TextField v-model="uploadAPIURL" required label="ATUS Tracker plugin URL"
            placeholder="e.g. https://your-tracker.to/atus/index.php" persistent-hint />

          <TextField v-model="uploadTrackerAnnounceURL" required label="Your trackers announce URL"
            placeholder="e.g. https://your-tracker.to/announce.php" persistent-hint
            hint="Your tracker announce URL without a passkey" />
        </v-card-text>
      </v-card>

      <v-card variant="text" title="Account Settings" class="card-accent mb-4">
        <v-card-text>
          <v-alert type="info" class="mb-4">
            Create a new user account on your tracker for the bot to use and enter the informations below.
            <br />
            The user will be visible as the uploader and seeder of releases on your tracker.
            <br />
            <i>The account does not need any special permissions.</i>
          </v-alert>

          <TextField v-model="uploadUserID" required label="User id" placeholder="e.g. 1" hint="" persistent-hint />

          <TextField v-model="uploadUserAnnounceURL" required label="The users announce URL with passkey"
            placeholder="e.g. https://your-tracker.to/announce.php?passkey=1234567890" persistent-hint
            hint="Your tracker announce URL WITH the users passkey" />
        </v-card-text>
      </v-card>

      <v-card variant="text" title="Other Settings" class="card-accent">
        <v-card-text>
          <TextField v-model="uploadCreatedBy" required label="Torrent created by" placeholder="e.g. ATUS"
            hint="Will be shown in some torrent clients" persistent-hint class="mb-2" />

          <TextField v-model="uploadComment" required label="Torrent upload comment"
            placeholder="e.g. my awesome tracker" persistent-hint hint="Will be shown in some torrent clients" />
        </v-card-text>
      </v-card>
    </v-card-text>
    <v-card-actions class="justify-end">
      <v-btn color="green" variant="tonal" :disabled="isLoading" @click="uploadTestTorrent()">
        Upload Test Release
      </v-btn>
      <v-btn color="primary" type="submit" :disabled="isLoading">Save</v-btn>
    </v-card-actions>
  </FormCard>
</template>



<script lang="ts">
import { defineComponent, ref } from "vue";
import { send } from "@/utils/websocket";
import { success, error } from "@/plugins/toast";

export default defineComponent({
  async setup() {
    const isLoading = ref(false);
    const form = ref(null);

    // --------------------------------------------------------------------------

    const uploadSuccessMessage = ref<any>(null);
    const uploadErrorMessage = ref<any>(null);

    const uploadUserID = ref("");
    const uploadUserAnnounceURL = ref("");
    const uploadTrackerAnnounceURL = ref("");
    const uploadAPIURL = ref("");
    const uploadCreatedBy = ref("");
    const uploadComment = ref("");

    // --------------------------------------------------------------------------

    const r: IResponse<IDestinationSettings> = await send("SETTINGS__UPLOAD__GET_ALL")
    uploadUserID.value = r.payload.uploadUserID;
    uploadUserAnnounceURL.value = r.payload.uploadUserAnnounceURL;
    uploadTrackerAnnounceURL.value = r.payload.uploadTrackerAnnounceURL;
    uploadAPIURL.value = r.payload.uploadAPIURL;
    uploadCreatedBy.value = r.payload.uploadCreatedBy;
    uploadComment.value = r.payload.uploadComment;

    // --------------------------------------------------------------------------

    const onSubmit = () => {
      isLoading.value = true;

      send("SETTINGS__UPLOAD__SAVE", {
        uploadUserID: uploadUserID.value,
        uploadUserAnnounceURL: uploadUserAnnounceURL.value,
        uploadTrackerAnnounceURL: uploadTrackerAnnounceURL.value,
        uploadAPIURL: uploadAPIURL.value,
        uploadCreatedBy: uploadCreatedBy.value,
        uploadComment: uploadComment.value,
      })
        .then(() => success("Settings saved successfully"))
        .catch(({ payload }: IResponse<string>) => error("Settings could not be saved", payload))
        .finally(() => isLoading.value = false);
    };

    // --------------------------------------------------------------------------

    const uploadTestTorrent = () => {
      isLoading.value = true;

      send(
        "SETTINGS__UPLOAD__UPLOAD_TEST_TORRENT", {
        uploadUserID: uploadUserID.value,
        uploadUserAnnounceURL: uploadUserAnnounceURL.value,
        uploadTrackerAnnounceURL: uploadTrackerAnnounceURL.value,
        uploadAPIURL: uploadAPIURL.value,
        uploadCreatedBy: uploadCreatedBy.value,
        uploadComment: uploadComment.value,
      })
        .then(({ payload }: IResponse<string>) => uploadSuccessMessage.value = payload)
        .catch(({ payload }: IResponse<string>) => uploadErrorMessage.value = payload)
        .finally(() => isLoading.value = false);
    };

    // --------------------------------------------------------------------------

    return {
      form,
      isLoading,
      uploadUserID,
      uploadUserAnnounceURL,
      uploadTrackerAnnounceURL,
      uploadAPIURL,
      uploadCreatedBy,
      uploadComment,
      onSubmit,
      uploadTestTorrent,
      uploadSuccessMessage,
      uploadErrorMessage,
    };
  },
});
</script>
