#!/usr/bin/env python
# asg.py - script to spin up/down all autoscaling instances.
# @TODO: select which ASG to spin up/down.
# @TODO: option to attach/detach and start/stop instances.
# @TODO: better interactive messages.
# @TODO: use Terraform not aws cli to maintain state.

import os
import sys
import boto3

REGION = 'us-east-1'

def main(argv):
  ec2 = boto3.resource('ec2')
  instances = ec2.instances.all()
  elb = boto3.client('elb')
  autoscaling = boto3.client('autoscaling')
  auto_scaling_groups = autoscaling.describe_auto_scaling_groups()['AutoScalingGroups']
  scheduled=[]

  if len(sys.argv) > 1 and argv[1] == 'up':
    print("Lights on")
    for i in instances:
      if i.state['Name'] == 'stopped':
        for t in i.tags:
          if t['Key'] == 'ScheduledUptime':
            if t['Value'] == 'True':
              response = ec2.Instance(i.id).start()
              print(' --> Starting Instance {}'.format(i.id))
              scheduled.append(i)
              break
      else:
        print(i.instance_id, " - instance is already running")
        #print(i.instance_id, " - instance is already running in autoscaling group:")

    # Register w/ load balancer
    for i in scheduled:
      for t in i.tags:
        if t['Key'] == 'ElasticLoadbalancerName':
          response = elb.register_instances_with_load_balancer(LoadBalancerName=t['Value'],Instances=[{'InstanceId': i.id}])
          print(" --> Registering instance {}".format(t['Value']))
          break

    # @FIXME: DesiredCapacity
    for asg in auto_scaling_groups:
      for t in asg['Tags']:
        if t['Key'] == 'ScheduledUptime':
          if t['Value'] == 'True':
            if asg['DesiredCapacity'] < 1:
              print(" --> Scaling up AutoScaling Group {}".format(asg['AutoScalingGroupName']))
              response = autoscaling.update_auto_scaling_group(AutoScalingGroupName=asg['AutoScalingGroupName'],DesiredCapacity=2,MinSize=0)
              break

  elif len(sys.argv) > 1 and argv[1] == 'down':
    for i in instances:
      #print(i, ":")
      get_instance_name = "aws ec2 describe-tags --region us-east-1 --filters 'Name=resource-id,Values={}' 'Name=key,Values=Name' --output text | cut -f5".format(i.instance_id)
      instance_name = os.system(get_instance_name)
      if i.state['Name'] == 'running':

        # Only stop instances assigned to an autoscaling group
        for t in i.tags:
          if t['Key'] == 'aws:autoscaling:groupName':
            print(" --> Detaching instance ", i.instance_id, "from autoscaling group and spinning down")
            # Only the autoscaling groups we want to scale down
            #if t['Value'] in asg_down:
            response = ec2.Instance(i.id).stop()
            print(' --> Stopping Instance {}'.format(i.id))
            #scheduled.append(i)
            break

    """
    for i in scheduled:
      for t in i.tags:
        if t['Key'] == 'ElasticLoadbalancerName':
          response = elb.deregister_instances_from_load_balancer(LoadBalancerName=t['Value'],Instances=[{'InstanceId': i.id}])
          print(' --> Deregistering Instance from load-balancer {}'.format(LoadBalancerName['Value']))
          break
    """
    for asg in auto_scaling_groups:
      print("Found autoscaling group:", asg['AutoScalingGroupName'])
      for t in asg['Tags']:
        if asg['DesiredCapacity'] > 0:
          print(" --> Scaling down AutoScaling Group {}".format(asg['AutoScalingGroupName']))
          response = autoscaling.update_auto_scaling_group(AutoScalingGroupName=asg['AutoScalingGroupName'],DesiredCapacity=0,MinSize=0)
          break
  else:
    print("Usage: asg up|down")

if __name__ == "__main__":
  main(sys.argv)
