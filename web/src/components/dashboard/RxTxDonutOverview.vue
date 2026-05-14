<script lang="ts" setup>
import { computed } from 'vue';
import { useTheme } from 'vuetify';
import { useI18n } from 'vue-i18n';
import type { RepositoryTotalBandwidths } from '@/api';
import { numberToFixer } from '@/utils/convertors';

const props = defineProps<{
    totalBandwidths: RepositoryTotalBandwidths;
}>();

const theme = useTheme();
const { t } = useI18n();

const chartOptions = computed(() => {
    const colors = theme.current.value.colors as Record<string, string>;
    const isDark = theme.global.current.value.dark;
    const foreColor = colors.textSecondary ?? colors['on-surface-variant'] ?? '#64748B';
    const emptySlice = colors.surface ?? colors.background ?? (isDark ? '#1E293B' : '#FFFFFF');

    return {
        labels: [t('TOTAL_TX'), t('TOTAL_RX')],
        chart: {
            type: 'donut',
            fontFamily: `inherit`,
            foreColor,
            toolbar: {
                show: false
            }
        },
        colors: [colors.primary, colors.lightprimary, emptySlice],
        plotOptions: {
            pie: {
                startAngle: 0,
                endAngle: 360,
                donut: {
                    size: '75%',
                    background: 'transparent'
                }
            }
        },
        stroke: {
            show: false
        },

        dataLabels: {
            enabled: false
        },
        legend: {
            show: false
        },
        tooltip: { theme: isDark ? 'dark' : 'light', fillSeriesColor: false }
    };
});
const txPercentage = computed(() => {
    const rx = props.totalBandwidths?.rx ?? 0;
    const tx = props.totalBandwidths?.tx ?? 0;
    const total = rx + tx;

    if (total === 0) return 0;
    return Math.round((tx / total) * 100);
});

const chart = computed(() => [+props.totalBandwidths?.tx.toFixed(6) || 0, +props.totalBandwidths?.rx.toFixed(6) || 0]);
</script>
<template>
    <v-card elevation="10">
        <v-card-item>
            <div class="d-sm-flex align-center justify-space-between pt-sm-2">
                <v-card-title class="text-h5 text-capitalize">TX / RX {{ t('OVERVIEW') }}</v-card-title>
            </div>
            <v-row>
                <v-col cols="7" sm="7">
                    <div class="mt-6">
                        <h6 class="text-h6 text-capitalize text-body-1">
                            {{ t('TOTAL') }} TX:
                            <br />
                            <span class="text-medium-emphasis"> {{ numberToFixer(props.totalBandwidths?.tx, 6) }} GB </span>
                        </h6>
                        <h6 class="text-h6 text-capitalize text-body-1">
                            {{ t('TOTAL') }} RX:
                            <br />
                            <span class="text-medium-emphasis text-body-1">
                                {{ numberToFixer(props.totalBandwidths?.rx, 6) }} GB
                            </span>
                        </h6>
                        <h6 class="text-h6 text-capitalize text-body-1 mt-5">
                            {{ t('AVERAGE') }} (TX):
                            <span class="text-medium-emphasis text-body-1"> {{ txPercentage }}% </span>
                        </h6>
                        <div class="d-flex align-center mt-sm-10 mt-8">
                            <h6 class="text-subtitle-1 text-medium-emphasis">
                                <v-icon class="mr-1" color="primary" icon="mdi mdi-checkbox-blank-circle" size="10" />
                                TX
                            </h6>
                            <h6 class="text-subtitle-1 text-medium-emphasis pl-5">
                                <v-icon
                                    class="mr-1"
                                    color="lightprimary"
                                    icon="mdi mdi-checkbox-blank-circle"
                                    size="10"
                                />
                                RX
                            </h6>
                        </div>
                    </div>
                </v-col>
                <v-col class="pl-lg-0" cols="5" sm="5">
                    <div class="d-flex align-center flex-shrink-0">
                        <apexchart :options="chartOptions" :series="chart" class="pt-6" height="145" type="donut" />
                    </div>
                </v-col>
            </v-row>
        </v-card-item>
    </v-card>
</template>
