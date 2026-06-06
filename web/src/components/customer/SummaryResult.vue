<script lang="ts" setup>
import type { CustomerIOSSetupResponse, CustomerSummaryResponse } from '@/api';
import { bytesToGB, bytesToTrafficSize, formatDate, trafficTypesTransformer } from '@/utils/convertors';
import UiChildCard from '@/components/shared/UiChildCard.vue';
import { useI18n } from 'vue-i18n';

defineProps<{
    result: CustomerSummaryResponse;
    iosSetup: CustomerIOSSetupResponse | null;
}>();

const emit = defineEmits(['newSummary', 'disconnect', 'downloadCertificate', 'loadIOSSetup']);

const { t } = useI18n();

const copyCertificatePassword = (password: string) => {
    navigator.clipboard?.writeText(password);
};
</script>

<template>
    <v-col cols="12" md="6">
        <UiChildCard class="px-3">
            <template #title-header>
                <v-row align="center" justify="start">
                    <v-col cols="12" md="8" sm="12">
                        <span class="text-capitalize text-primary text-h3">
                            <span>{{ result.ocserv_user.username }}</span>
                            <span class="text-muted mx-1 text-capitalize"> ({{ t('ACCOUNT_AND_USAGE_SUMMARY') }}) </span>
                        </span>
                    </v-col>
                    <v-col cols="12" md="auto" sm="12">
                        <v-btn color="primary" flat size="small" @click="emit('disconnect')">
                            {{ t('DISCONNECT_ALL_SESSIONS') }}
                        </v-btn>
                    </v-col>
                    <v-col v-if="result.ocserv_user.certificate_available" cols="12" md="auto" sm="12">
                        <v-btn color="info" size="small" variant="outlined" @click="emit('downloadCertificate')">
                            {{ t('DOWNLOAD_CERTIFICATE') }}
                        </v-btn>
                    </v-col>
                    <v-col v-if="result.ocserv_user.certificate_available" cols="12" md="auto" sm="12">
                        <v-btn color="primary" size="small" variant="outlined" @click="emit('loadIOSSetup')">
                            {{ t('IOS_SETUP') }}
                        </v-btn>
                    </v-col>
                </v-row>
            </template>

            <template #action>
                <v-spacer />
                <v-btn color="primary" @click="emit('newSummary')">
                    {{ t('CHECK_ANOTHER') }}
                </v-btn>
            </template>

            <div class="space-y-4 mt-8 px-1">
                <!-- General info -->
                <div class="bg-surface shadow rounded-lg p-4">
                    <h4 class="text-lg font-semibold my-4">{{ t('DETAILS') }}</h4>

                    <div class="grid grid-cols-2 gap-4 mx-5">
                        <v-row align="center" justify="start">
                            <v-col cols="12" md="6">
                                <span class="font-medium text-gray-600 text-capitalize"> {{ t('TRAFFIC_TYPE') }}: </span>
                                <span class="ms-1 text-capitalize">
                                    {{ trafficTypesTransformer(result.ocserv_user.traffic_type) }}
                                </span>
                            </v-col>

                            <v-col cols="12" md="6">
                                <span class="font-medium text-gray-600 text-capitalize"> {{ t('TRAFFIC_SIZE') }}: </span>
                                <span class="ms-1">{{ bytesToTrafficSize(result.ocserv_user.traffic_size) }}</span>
                            </v-col>
                        </v-row>
                        <v-row align="center" justify="start">
                            <v-col cols="12" md="6">
                                <span class="font-medium text-gray-600 text-capitalize"> {{ t('CERTIFICATE') }}: </span>
                                <span
                                    :class="result.ocserv_user.certificate_enabled ? 'text-success' : 'text-warning'"
                                    class="ms-1 text-capitalize"
                                >
                                    {{ result.ocserv_user.certificate_enabled ? t('ENABLED') : t('DISABLED') }}
                                </span>
                            </v-col>

                            <v-col cols="12" md="6">
                                <span class="font-medium text-gray-600 text-capitalize"> {{ t('TOTAL') }}: </span>
                                <span class="ms-1">
                                    {{ bytesToGB(result.ocserv_user.rx + result.ocserv_user.tx) }} GB
                                </span>
                            </v-col>
                        </v-row>
                        <v-row align="center" justify="start">
                            <v-col cols="12" md="6">
                                <span class="font-medium text-gray-600 text-capitalize"> RX ({{ t('TOTAL') }}): </span>
                                <span class="ms-1">{{ bytesToGB(result.ocserv_user.rx) }} GB</span>
                            </v-col>

                            <v-col cols="12" md="6">
                                <span class="font-medium text-gray-600 text-capitalize"> TX ({{ t('TOTAL') }}): </span>
                                <span class="ms-1">{{ bytesToGB(result.ocserv_user.tx) }} GB</span>
                            </v-col>
                        </v-row>
                        <v-row align="center" justify="start">
                            <v-col cols="12" md="6">
                                <span class="font-medium text-gray-600 text-capitalize"> {{ t('EXPIRE_AT') }}: </span>
                                <span v-if="result.ocserv_user.expire_at" class="ms-1">
                                    {{ formatDate(result.ocserv_user.expire_at) }}
                                </span>
                                <span v-else class="ms-1 text-warning italic">{{ t('NOT_SET') }}</span>
                            </v-col>

                            <v-col cols="12" md="6">
                                <span class="font-medium text-gray-600 text-capitalize">
                                    {{ t('DEACTIVATED_AT') }}:
                                </span>
                                <span v-if="result.ocserv_user.deactivated_at" class="ms-1">
                                    {{ formatDate(result.ocserv_user.deactivated_at) }}
                                </span>
                                <span v-else class="ms-1 text-warning italic">{{ t('NOT_SET') }}</span>
                            </v-col>
                        </v-row>
                    </div>

                    <h4 class="text-lg font-semibold my-4">{{ t('MONTHLY_BANDWIDTHS') }}</h4>

                    <div class="grid grid-cols-2 gap-4 mx-5">
                        <v-row align="center" justify="start">
                            <v-col cols="12" md="6">
                                <span class="font-medium text-gray-600 text-capitalize"> {{ t('DATE_START') }}: </span>
                                <span class="ms-1">{{ formatDate(result.usage.date_start) }}</span>
                            </v-col>
                            <v-col cols="12" md="6">
                                <span class="font-medium text-gray-600 text-capitalize"> {{ t('DATE_END') }}: </span>
                                <span class="ms-1">{{ formatDate(result.usage.date_end) }}</span>
                            </v-col>
                        </v-row>
                        <v-row align="center" justify="start">
                            <v-col cols="12" md="6">
                                <span class="font-medium text-gray-600 text-capitalize"> RX: </span>
                                <span class="ms-1">{{ bytesToGB(result.usage.bandwidths.rx) }} GB</span>
                            </v-col>

                            <v-col cols="12" md="6">
                                <span class="font-medium text-gray-600 text-capitalize"> TX: </span>
                                <span class="ms-1">{{ bytesToGB(result.usage.bandwidths.tx) }} GB</span>
                            </v-col>
                        </v-row>
                    </div>
                </div>
            </div>
        </UiChildCard>

        <UiChildCard v-if="iosSetup" class="px-3 mt-4">
            <template #title-header>
                <span class="text-capitalize text-primary text-h3">
                    {{ t('IOS_SETUP_TITLE') }}
                </span>
            </template>

            <div class="space-y-4 mt-8 px-1">
                <v-alert class="mb-4" type="info" variant="tonal">
                    {{ t('IOS_SETUP_EXTERNAL_CONTROL_HINT') }}
                </v-alert>

                <v-alert class="mb-4" type="warning" variant="tonal">
                    {{ t('IOS_SETUP_PASSWORD_HINT') }}
                    <strong>{{ iosSetup.certificate_password }}</strong>
                    <v-btn
                        class="ms-2"
                        size="x-small"
                        variant="text"
                        @click="copyCertificatePassword(iosSetup.certificate_password)"
                    >
                        {{ t('COPY') }}
                    </v-btn>
                </v-alert>

                <v-row>
                    <v-col cols="12" md="6">
                        <v-btn block color="primary" :href="iosSetup.certificate_import_uri" variant="flat">
                            {{ t('IOS_IMPORT_CERTIFICATE') }}
                        </v-btn>
                    </v-col>

                    <v-col cols="12" md="6">
                        <v-btn block color="success" :href="iosSetup.connection_create_uri" variant="flat">
                            {{ t('IOS_ADD_CONNECTION') }}
                        </v-btn>
                    </v-col>
                </v-row>

                <p class="mt-4 text-caption">
                    {{ t('IOS_SETUP_EXPIRES_HINT') }}
                </p>
            </div>
        </UiChildCard>
    </v-col>
</template>
