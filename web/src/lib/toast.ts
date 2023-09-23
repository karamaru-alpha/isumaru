import { toast } from '@zerodevx/svelte-toast'

export const success = (message: string): void => {
    toast.push(message, {
        theme: {
            "--toastBackground": "rgba(0,200,100,0.8)",
            "--toastBarBackground": "rgba(0,200,100,0.9)",
        }
    });
};

export const error = (message: string): void => {
    toast.push(message, {
        theme: {
            "--toastBackground": "rgba(255,0,0,0.7)",
            "--toastBarBackground": "rgba(255,0,0,0.8)"
        }
    });
};
