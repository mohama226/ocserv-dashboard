<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { onBeforeUnmount, onMounted, ref, watch, type PropType } from 'vue';

type CurrentAction = {
    value: string | null;
    expiresAt: number;
};

const props = defineProps({
    state: {
        type: String as PropType<'active' | 'inactive' | 'failed' | 'activating' | 'deactivating'>,
        default: 'inactive'
    }
});

const emit = defineEmits(['getState']);
const { t } = useI18n();

const currentAction = ref<string | null>(null);
const remainingTime = ref<number>(0);
const storageKey = 'current_action';
const TTL = 2 * 60 * 1000; // 2 minutes

let validateInterval: ReturnType<typeof setInterval> | null = null;
let countdownInterval: ReturnType<typeof setInterval> | null = null;
let hasExpired = false;

const validateStorage = () => {
    const raw = localStorage.getItem(storageKey);

    if (!raw) {
        currentAction.value = null;
        remainingTime.value = 0;
        return;
    }

    try {
        const data: CurrentAction = JSON.parse(raw);

        const diff = data.expiresAt - Date.now();

        if (diff <= 0) {
            localStorage.removeItem(storageKey);
            currentAction.value = null;
            remainingTime.value = 0;

            if (!hasExpired) {
                hasExpired = true;
                emit('getState');
            }

            return;
        }

        currentAction.value = data.value;
        remainingTime.value = Math.ceil(diff / 1000);
    } catch {
        localStorage.removeItem(storageKey);
        currentAction.value = null;
        remainingTime.value = 0;
    }
};

const setCurrentAction = (value: string) => {
    const payload: CurrentAction = {
        value,
        expiresAt: Date.now() + TTL
    };

    localStorage.setItem(storageKey, JSON.stringify(payload));
    currentAction.value = value;
    hasExpired = false;
};

const status = () => {
    emit('getState');
};

const restart = () => {
    if (props.state !== 'active') return;

    emit('getState');
    setCurrentAction('restart');
};

const enable = () => {
    if (props.state !== 'inactive') return;

    emit('getState');
    setCurrentAction('enable');
};

const disable = () => {
    if (props.state !== 'active') return;

    emit('getState');
    setCurrentAction('disable');
};

onMounted(() => {
    validateStorage();
});

watch(
    currentAction,
    (val) => {
        // reset intervals
        if (validateInterval) clearInterval(validateInterval);
        if (countdownInterval) clearInterval(countdownInterval);

        validateInterval = null;
        countdownInterval = null;

        if (!val) {
            hasExpired = false;
            return;
        }

        // 30s validation
        validateInterval = setInterval(() => {
            validateStorage();
        }, 30000);

        // 1s countdown
        countdownInterval = setInterval(() => {
            const raw = localStorage.getItem(storageKey);

            if (!raw) {
                currentAction.value = null;
                remainingTime.value = 0;
                return;
            }

            try {
                const data: CurrentAction = JSON.parse(raw);
                const diff = data.expiresAt - Date.now();

                if (diff <= 0) {
                    localStorage.removeItem(storageKey);
                    currentAction.value = null;
                    remainingTime.value = 0;

                    if (!hasExpired) {
                        hasExpired = true;
                        emit('getState');
                    }

                    return;
                }

                remainingTime.value = Math.ceil(diff / 1000);
            } catch {
                localStorage.removeItem(storageKey);
                currentAction.value = null;
                remainingTime.value = 0;
            }
        }, 1000);
    },
    { immediate: true }
);

onBeforeUnmount(() => {
    if (validateInterval) clearInterval(validateInterval);
    if (countdownInterval) clearInterval(countdownInterval);
});
</script>

<template>
    <v-row align="center" justify="center" v-if="currentAction == null">
        <v-col cols="12" md="auto" v-if="props.state == 'active'">
            <v-btn @click="status" color="info"> {{ t('RELOAD') }} {{ t('STATUS') }} </v-btn>
        </v-col>

        <v-col cols="12" md="auto" v-if="props.state == 'active'">
            <v-btn @click="restart" color="primary">
                {{ t('RESTART') }}
            </v-btn>
        </v-col>

        <v-col cols="12" md="auto" v-if="state == 'inactive'">
            <v-btn @click="enable" color="success">
                {{ t('ENABLE') }}
            </v-btn>
        </v-col>

        <v-col cols="12" md="auto" v-if="state == 'active'">
            <v-btn @click="disable" color="error">
                {{ t('DISABLE') }}
            </v-btn>
        </v-col>
    </v-row>

    <div v-else class="text-center">
        <div class="text-h6 text-primary text-capitalize">{{ currentAction }}ing ...</div>
        <div v-if="remainingTime > 0" class="text-h6 my-3" style="color: #888888">
            {{ t('PLEASE_WAIT_UNTIL') }}: {{ remainingTime }}s
        </div>
    </div>
</template>
