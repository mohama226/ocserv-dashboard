<script lang="ts" setup>
import { useI18n } from 'vue-i18n';
import { SystemdApi, type SystemdOcservSystemdStatus } from '@/api';
import { getAuthorization } from '@/utils/request';
import { computed, onMounted, ref } from 'vue';

const emit = defineEmits(['state']);

const { t } = useI18n();

const service = ref<SystemdOcservSystemdStatus>({});

const getStatus = () => {
    const api = new SystemdApi();

    api.systemdStatusGet({
        ...getAuthorization()
    }).then((res) => {
        Object.assign(service.value, res.data);
        emit('state', res.data.active_state);
    });
};

onMounted(() => {
    getStatus();
});

defineExpose({
    getStatus
});

// status color
const statusColor = computed(() => {
    const active = service.value.active_state;
    const unit = service.value.unit_file_state;

    if (active === 'active') return 'success';
    if (active === 'activating') return 'info';
    if (active === 'deactivating') return 'warning';
    if (active === 'failed') return 'error';
    if (active === 'inactive' && unit === 'disabled') return 'grey';
    if (unit === 'masked') return 'black';

    return 'secondary';
});

const formatMemory = (bytes?: number) => {
    if (!bytes) return '-';

    const mb = bytes / 1024 / 1024;
    return `${mb.toFixed(1)} MB`;
};

const formatCPU = (ns?: number) => {
    if (!ns) return '-';

    const sec = ns / 1e9;
    return `${sec.toFixed(2)} s`;
};
</script>

<template>
    <v-row>
        <v-col cols="12" md="12">
            <!-- HEADER -->
            <v-row align="center" justify="space-between" class="mt-2 text-capitalize">
                <v-col class="ma-0 pa-0 ms-6">
                    {{ service.description }}
                </v-col>
                <v-col class="ma-0 pa-0 me-5">
                    <v-chip :color="statusColor">
                        {{ service.active_state }}
                    </v-chip>
                </v-col>
            </v-row>

            <v-divider class="my-3" />

            <!-- GRID INFO -->
            <v-row dense>
                <v-col cols="12" md="6">
                    <v-list>
                        <v-list-item class="mb-3">
                            <v-list-item-title>ID</v-list-item-title>
                            <v-list-item-subtitle>{{ service.id }}</v-list-item-subtitle>
                        </v-list-item>

                        <v-list-item class="mb-3">
                            <v-list-item-title>{{ t('SUB_STATE') }}</v-list-item-title>
                            <v-list-item-subtitle>{{ service.sub_state }}</v-list-item-subtitle>
                        </v-list-item>

                        <v-list-item class="mb-3">
                            <v-list-item-title>{{ t('UNIT_STATE') }}</v-list-item-title>
                            <v-list-item-subtitle>{{ service.unit_file_state }}</v-list-item-subtitle>
                        </v-list-item>

                        <v-list-item class="mb-3">
                            <v-list-item-title>{{ t('START_TIME') }}</v-list-item-title>
                            <v-list-item-subtitle>{{ service.start_time }}</v-list-item-subtitle>
                        </v-list-item>
                    </v-list>
                </v-col>

                <v-col cols="12" md="6">
                    <v-list>
                        <v-list-item class="mb-3">
                            <v-list-item-title>PID</v-list-item-title>
                            <v-list-item-subtitle>{{ service.main_pid }}</v-list-item-subtitle>
                        </v-list-item>

                        <v-list-item class="mb-3">
                            <v-list-item-title>{{ t('MEMORY') }}</v-list-item-title>
                            <v-list-item-subtitle>
                                {{ formatMemory(service.memory) }}
                            </v-list-item-subtitle>
                        </v-list-item>

                        <v-list-item class="mb-3">
                            <v-list-item-title>{{ t('CPU_USAGE') }}</v-list-item-title>
                            <v-list-item-subtitle>
                                {{ formatCPU(service.cpu_usage_nsec) }}
                            </v-list-item-subtitle>
                        </v-list-item>

                        <v-list-item class="mb-3">
                            <v-list-item-title>{{ t('TASKS') }}</v-list-item-title>
                            <v-list-item-subtitle>{{ service.tasks }}</v-list-item-subtitle>
                        </v-list-item>
                    </v-list>
                </v-col>
            </v-row>
        </v-col>
    </v-row>
</template>
