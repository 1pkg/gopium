export default interface Confs {
    args: Args,
    flags: Flags,
}

export interface Args {
    package: string;
    walker: string;
    startegies: string[];
}

export interface Flags {
    // target platform vars
    c?: string;
    a?: string;
    l?: number[];
    // package parser vars
    p?: string;
    e?: string[];
    f?: string[];
    // walker vars
    r?: string
    d?: boolean;
    b?: boolean;
    // printer vars
    i?: number;
    w?: number;
    s?: boolean;
    // global vars
    t?: number;
}